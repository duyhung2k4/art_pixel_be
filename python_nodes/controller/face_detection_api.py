from flask import Blueprint, request, jsonify
import cv2
import face_recognition
import os

face_detection_bp = Blueprint('face_detection', __name__)

@face_detection_bp.route('/detect_single_face', methods=['POST'])
def detect_single_face():
    image_path = request.json.get("input_image_path")
    
    if not image_path or not os.path.exists(image_path):
        return jsonify({"result": False, "error": "Image path is invalid."}), 400

    try:
        image = cv2.imread(image_path)
        face_locations = face_recognition.face_locations(image)
        is_single_face = len(face_locations) == 1
        
        if is_single_face == True:
            return jsonify({"result": True})

        return jsonify({"result": False})
    except Exception as e:
        return jsonify({"result": False, "error": str(e)}), 500
