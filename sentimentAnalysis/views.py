# ? VIEWS.PY

# Importing dependencies.
from flask import Flask, Blueprint, render_template, request, jsonify
from transformers import AutoTokenizer, AutoModelForSequenceClassification
from scipy.special import softmax
from goose3 import Goose
from werkzeug.utils import secure_filename
import numpy as np
import os
import mimetypes

# Creating a blueprint for views to use for routing.
app = Flask(__name__)
views = Blueprint(__name__, "views")

HUGGINGFACE_CACHE_DIR = "./staticFileCache"
TORCH_CACHE_DIR = "./staticFileCache"
os.environ["TORCH_HOME"] = TORCH_CACHE_DIR

# Importing pre-trained model for Sentiment Analysis.
SENTIMENT_MODEL = "cardiffnlp/twitter-roberta-base-sentiment-latest"

# Model Training for Polarity Scores.
tokenizer = AutoTokenizer.from_pretrained(
    SENTIMENT_MODEL, cache_dir=HUGGINGFACE_CACHE_DIR
)
sentiment_model = AutoModelForSequenceClassification.from_pretrained(
    SENTIMENT_MODEL, cache_dir=HUGGINGFACE_CACHE_DIR
)


@views.route("/analysis", methods=["POST"])
def home():

    # When a form input is received, show the sentiment based on the input.
    data = request.json

    input_type = data["type"]

    # Find the input text from different types.
    input_text = ""

    print(data["value"])

    if input_type == "text":
        input_text = data["value"]
    elif input_type == "url":
        url = request.form.get("input")
        g = Goose()
        article = g.extract(url=url)
        input_text = article.cleaned_text
    elif input_type == "media":
        file = request.files.get("input")
        if file:
            filename = secure_filename(file.filename)
            file_path = os.path.join(app.root_path, "static", "files", filename)
            file.save(file_path)
            input_media = process_files(file_path)

    # Find the sentiment values.
    sentiment_analysis = find_text_sentiment_analysis(input_text)
    print(jsonify(sentiment_analysis))

    return jsonify(sentiment_analysis)


# * Chunk the text into pieces of 510 characters.
def chunk_text(text, max_len=510):
    sentences = text.split(". ")
    chunks = []
    current_chunk = ""

    for sentence in sentences:
        # If the chunk is less than max length.
        if len(current_chunk) + len(sentence) < max_len:
            if current_chunk:
                # Add a space for continuing sentences.
                current_chunk += " "
            current_chunk += sentence
        # If the chunk is more than max length.
        else:
            chunks.append(current_chunk)
            current_chunk = sentence

        # Adding the last chunk to chunks.
        chunks.append(current_chunk)

        return chunks


# * Process Media Files for analysis.
def process_files(file_path):
    mime_type, encoding = mimetypes.guess_type(file_path)
    type, subtype = mime_type.split("/", 1)


# * Find the polarity scores of the input.
def find_text_sentiment_analysis(input):

    # Split the input into separate chunks.
    chunks = chunk_text(input)
    sentiment_dicts = []

    for chunk in chunks:
        # Find tokenized words.
        encoded_text = tokenizer(chunk, return_tensors="pt")

        # Find polarity scores.
        output = sentiment_model(**encoded_text)
        scores = output[0][0].detach().numpy()
        scores = softmax(scores)

        # Scores
        val_neg = scores[0]
        val_neu = scores[1]
        val_pos = scores[2]

        # Find Prominent Sentiment
        if (val_neg > val_pos) and (val_neg > val_neu):
            prominent_sentiment = "NEGATIVE"
        elif (val_pos > val_neg) and (val_pos > val_neu):
            prominent_sentiment = "POSITIVE"
        else:
            prominent_sentiment = "NEUTRAL"

        # Create Sentiment Analysis Dictionary
        sentiment_dict = {
            "score_negative": float(val_neg),
            "score_neutral":  float(val_neu),
            "score_positive":  float(val_pos),
            "prominent_sentiment": prominent_sentiment,
        }

        sentiment_dicts.append(sentiment_dict)

    # Aggregate the list of chunks to find average sentiment.
    avg_sentiment_dict = {
        "score_negative":
            np.mean([float(d["score_negative"]) for d in sentiment_dicts])
    ,
        "score_neutral": 
            np.mean([float(d["score_neutral"]) for d in sentiment_dicts])
        ,
        "score_positive": 
            np.mean([float(d["score_positive"]) for d in sentiment_dicts])
        ,
        "prominent_sentiment": max(
            set([d["prominent_sentiment"] for d in sentiment_dicts]),
            key=[d["prominent_sentiment"] for d in sentiment_dicts].count,
        ),
    }

    return avg_sentiment_dict
