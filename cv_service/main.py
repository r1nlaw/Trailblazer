import os
import torch
from sklearn.metrics.pairwise import cosine_similarity
from PIL import Image
from transformers import CLIPProcessor, CLIPModel
import numpy as np

# Получаем текущую рабочую директорию
current_directory = os.getcwd()

# Формируем пути относительно текущей директории
query_folder = os.path.join(current_directory, "input_images")  # Папка с изображениями для сравнения
embeddings_folder = os.path.join(current_directory, "references")  # Папка с предвычисленными эмбеддингами

# Выбор устройства
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

# Загрузка модели CLIP
model = CLIPModel.from_pretrained("openai/clip-vit-large-patch14").to(device)
processor = CLIPProcessor.from_pretrained("openai/clip-vit-large-patch14", use_fast=True)

# Функция для получения эмбеддинга изображения
def get_clip_embedding(image_path):
    image = Image.open(image_path).convert("RGB")
    inputs = processor(images=image, return_tensors="pt").to(device)

    with torch.no_grad():
        embedding = model.get_image_features(**inputs)
        embedding = embedding / embedding.norm(dim=-1, keepdim=True)  # Нормализация
    return embedding.cpu().numpy()

# Функция для загрузки эмбеддингов из файлов
def load_embeddings_from_files(embeddings_folder):
    embeddings = {}
    for file_name in os.listdir(embeddings_folder):
        file_path = os.path.join(embeddings_folder, file_name)
        if file_path.endswith('.npy'):  # Проверка на формат .npy для эмбеддингов
            embedding = np.load(file_path)
            embeddings[file_name] = embedding
    return embeddings

# Функция сравнения изображения с эмбеддингами
def compare_images_with_embeddings(query_folder, embeddings_folder, similarity_threshold=0.80):
    # Загружаем эмбеддинги из папки
    embeddings = load_embeddings_from_files(embeddings_folder)

    # Перебираем все изображения в query_folder
    for query_image_name in os.listdir(query_folder):
        query_image_path = os.path.join(query_folder, query_image_name)

        if query_image_path.lower().endswith(('.jpg', '.jpeg', '.png')):

            print(f"\nСравнение для изображения: {query_image_name}")
            
            # Получаем эмбеддинг для изображения
            query_embedding = get_clip_embedding(query_image_path)

            # Сравниваем с эмбеддингами из папки
            similarities = []
            for ref_name, ref_embedding in embeddings.items():
                similarity = cosine_similarity(query_embedding, ref_embedding)[0][0]
                similarities.append((ref_name, similarity))

            # Сортируем по схожести
            similarities.sort(key=lambda x: x[1], reverse=True)

            # Выводим топ-3 самых схожих изображений
            print("Топ-5 схожих изображений:")

            # Проверяем, все ли топ-3 изображений имеют схожесть выше порога
            all_similar = True
            for image_name, similarity in similarities[:3]:
                print(f"{image_name}: {similarity * 100:.2f}% схожести")
                if similarity < similarity_threshold:
                    all_similar = False

            # Если все 3 изображений имеют схожесть выше порога, считаем изображение схожим
            if all_similar:
                print("  --> СОВПАДАЕТ (все 3 изображения схожи на 80% и выше)")
            else:
                print("  --> НЕ СОВПАДАЕТ (не все изображения в топ-5 схожи на 80%)")

# Выполнение функции сравнения
compare_images_with_embeddings(query_folder, embeddings_folder)
