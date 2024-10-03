from flask import Blueprint, request, jsonify
from deepface import DeepFace
import os
import numpy as np

face_encoding_bp = Blueprint('face_encoding', __name__)

@face_encoding_bp.route('/face_encoding', methods=['POST'])
def face_encoding():
    directory_path = request.json.get("directory_path")

    if not directory_path or not os.path.exists(directory_path):
        return jsonify({"result": "error", "message": "Directory path is invalid."}), 400

    list_face_encoding = []

    try:
        for image_file in os.listdir(directory_path):
            image_path = os.path.join(directory_path, image_file)
            
            print(image_path)

            # Mã hóa khuôn mặt sử dụng DeepFace
            embedding = DeepFace.represent(img_path=image_path, model_name="VGG-Face", enforce_detection=False)
            
            if embedding:
                list_face_encoding.append(embedding)

        return jsonify({"result": "success", "face_encodings": list_face_encoding})
    
    except Exception as e:
        return jsonify({"result": "error", "message": str(e)}), 500
