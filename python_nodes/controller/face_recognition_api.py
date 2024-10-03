from flask import Blueprint, request, jsonify
import json
import face_recognition
import numpy as np
import os

face_recognition_bp = Blueprint('face_recognition', __name__)

@face_recognition_bp.route('/recognize_faces', methods=['POST'])
def recognize_faces_from_db():
    data = request.json

    try:
        faces = data["faces"]
        known_face_encodings = [np.array(face["faceEncoding"]) for face in faces]
        known_profile_ids = [face["profileId"] for face in faces]

        input_image_path = data["input_image_path"]
        if not os.path.exists(input_image_path):
            return jsonify({"result": "-2", "message": "Input image path does not exist."}), 400

    except KeyError as e:
        return jsonify({"result": "-2", "error": f"Missing key: {str(e)}"}), 400
    except Exception as e:
        return jsonify({"result": "-2", "error": str(e)}), 400

    try:
        def recognize_face_in_image(input_image_path):
            image_to_check = face_recognition.load_image_file(input_image_path)
            face_locations = face_recognition.face_locations(image_to_check)
            face_encodings = face_recognition.face_encodings(image_to_check, face_locations)

            if len(face_encodings) == 0:
                return -3  # No faces detected

            for face_encoding in face_encodings:
                matches = face_recognition.compare_faces(known_face_encodings, face_encoding)

                if not any(matches):
                    continue  # No match found, continue to next face encoding

                face_distances = face_recognition.face_distance(known_face_encodings, face_encoding)
                best_match_index = np.argmin(face_distances)

                if matches[best_match_index]:
                    profile_id = known_profile_ids[best_match_index]
                    accuracy = 1 - face_distances[best_match_index]  # Calculate accuracy
                    
                    print(f"ProfileID:{profile_id} / result: {round(accuracy * 100, 2)}")

                    if round(accuracy * 100, 2) >= 80.00:
                        return profile_id
                    return -3  # Return only the profile ID

            return -3  # No match found

        profile_id = recognize_face_in_image(input_image_path)
        if profile_id == -3:
            return jsonify({"result": "-3", "message": "No matching faces found."})
        
        return jsonify({"result": str(profile_id)})

    except Exception as e:
        return jsonify({"result": "-4", "error": str(e)}), 500
