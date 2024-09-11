import sys
import cv2
import face_recognition

def detect_single_face(image_path):
    image = cv2.imread(image_path)
    face_locations = face_recognition.face_locations(image)
    return len(face_locations) == 1

if __name__ == "__main__":
    image_path = sys.argv[1]  # Nhận đường dẫn từ tham số dòng lệnh
    result = detect_single_face(image_path)
    if result:
        print("true")
    else:
        print("false")
