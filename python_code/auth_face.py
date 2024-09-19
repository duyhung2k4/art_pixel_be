import sys
import json
import face_recognition
import numpy as np

def recognize_faces_from_db(json_file_path):
    try:
        # Mở file JSON và đọc nội dung
        with open(json_file_path, 'r') as f:
            face_data = f.read()
            data = json.loads(face_data)  # Chuyển JSON thành dictionary
    except Exception as e:
        return "-3"

    try:
        # Tách mã hóa khuôn mặt và ProfileId từ phần "faces"
        faces = data["faces"]
        known_face_encodings = [np.array(face["faceEncoding"]) for face in faces]
        known_profile_ids = [face["profileId"] for face in faces]  # Chỉ dùng ProfileId

        # Lấy đường dẫn tới ảnh từ input_image_path
        input_image_path = data["input_image_path"]
        # Kiểm tra xem file ảnh có tồn tại không
        with open(input_image_path, 'rb') as f:
            pass
    except Exception as e:
        return "-2"

    try:
        # Hàm nhận diện khuôn mặt từ hình ảnh đầu vào
        def recognize_face_in_image(input_image_path):
            # Tải ảnh cần nhận diện
            image_to_check = face_recognition.load_image_file(input_image_path)

            # Tìm vị trí và mã hóa các khuôn mặt trong ảnh
            face_locations = face_recognition.face_locations(image_to_check)
            face_encodings = face_recognition.face_encodings(image_to_check, face_locations)

            if len(face_locations) == 0:
                return "-1"

            for face_encoding in face_encodings:
                # So sánh với các khuôn mặt đã huấn luyện
                matches = face_recognition.compare_faces(known_face_encodings, face_encoding)
                
                # Tìm khoảng cách giữa khuôn mặt được phát hiện và các khuôn mặt đã huấn luyện
                face_distances = face_recognition.face_distance(known_face_encodings, face_encoding)
                best_match_index = np.argmin(face_distances)
                
                # Kiểm tra xem có khớp với khuôn mặt đã biết không
                if matches[best_match_index]:
                    profile_id = known_profile_ids[best_match_index]
                    
                    # Kiểm tra đơn giản: ví dụ, kiểm tra số lượng khuôn mặt phát hiện
                    if len(face_locations) == 1:  # Chỉ chấp nhận 1 khuôn mặt duy nhất
                        return f"{profile_id}"

            return "-1"  # Nếu không tìm thấy sự khớp

        # Gọi hàm nhận diện khuôn mặt với đường dẫn ảnh
        message = recognize_face_in_image(input_image_path)
        return message
    except Exception as e:
        return "-4"

if __name__ == "__main__":
    # Nhận đường dẫn file JSON từ Go qua argv
    json_file_path = sys.argv[1]
    result = recognize_faces_from_db(json_file_path)
    print(result)
