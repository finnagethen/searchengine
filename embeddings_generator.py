from sentence_transformers import SentenceTransformer
import numpy as np
import json
import struct
import argparse


class Embedder:
    """
    A class to generate and export document embeddings using a specified
    SentenceTransformer model.
    """

    def __init__(self, model_name: str):
        self.model_name = model_name
        self.model = SentenceTransformer(model_name)
        self.embeddings: np.ndarray | None = None

    def build_from_file(self, path: str):
        """
        Reads documents from a file and generates embeddings for each document
        using the specified model.
        The file is expected to contain one document per line.
        """
        documents: list[str] = []

        with open(path, "r", encoding="utf-8") as file:
            # Skip header
            next(file, None)

            for document in file:
                documents.append(document)

        self.embeddings = self.model.encode(
            documents,
            normalize_embeddings=True,
            show_progress_bar=True,
            convert_to_numpy=True,
        )

        # Ensure the embeddings are in float32 format for efficient storage.
        self.embeddings = self.embeddings.astype(np.float32)

    def export(self, embeddings_path: str, metadata_path: str):
        """
        Exports the generated embeddings to a binary file and saves metadata to
        a JSON file.
        """
        if self.embeddings is None:
            raise ValueError("Embeddings have not been generated yet.")

        # Save embeddings to a binary file.
        # The binary file has a header containing the number of documents and
        # the dimension of the embeddings, followed by the raw embedding data.
        with open(embeddings_path, "wb") as file:
            # uint32 num_docs
            # uint32 dimension
            num_docs, dimension = self.embeddings.shape
            file.write(struct.pack("<I", num_docs))
            file.write(struct.pack("<I", dimension))

            # Raw float32 embedding data
            self.embeddings.tofile(file)

        # Save metadata to a JSON file.
        # The metadata includes the model, dimension, whether the embeddings
        # are normalized, the data type, and the number of documents.
        metadata = {
            "model": self.model_name,
            "dimension": dimension,
            "normalized": True,
            "dtype": "float32",
            "num_docs": num_docs,
        }

        with open(metadata_path, "w", encoding="utf-8") as f:
            json.dump(metadata, f, indent=4)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str, help="file containing the documents to embed")
    parser.add_argument(
        "--model",
        type=str,
        default="all-MiniLM-L6-v2",
        help="the name of the SentenceTransformer model to use",
    )
    parser.add_argument(
        "--output-embeddings",
        type=str,
        default="embeddings.bin",
        help="the name of the output file to write the embeddings to",
    )
    parser.add_argument(
        "--output-metadata",
        type=str,
        default="metadata.json",
        help="the name of the output file to write the metadata to",
    )
    return parser.parse_args()


def main(args: argparse.Namespace):
    """
    Generates document embeddings using the specified model and saves them to a binary file.
    """

    print(f"Loading model `{args.model}` ...")
    embedder = Embedder(args.model)

    print(f"Building embeddings from file `{args.file}` ...")
    embedder.build_from_file(args.file)

    print(
        f"Exporting embeddings to `{args.output_embeddings}` and metadata to `{args.output_metadata}` ..."
    )
    embedder.export(args.output_embeddings, args.output_metadata)

    print("Done.")


if __name__ == "__main__":
    main(parse_args())

