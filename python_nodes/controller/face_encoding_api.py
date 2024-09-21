from flask import Blueprint, request, jsonify
import face_recognition
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
            new_image = face_recognition.load_image_file(image_path)
            new_face_encodings = face_recognition.face_encodings(new_image)
            
            if len(new_face_encodings) > 0:
                new_face_encoding = new_face_encodings[0]
                list_face_encoding.append(new_face_encoding.tolist())
            else:
                print(f"No faces found in {image_file}")

        return jsonify({"result": "success", "face_encodings": list_face_encoding})

    except Exception as e:
        return jsonify({"result": "error", "message": str(e)}), 500
