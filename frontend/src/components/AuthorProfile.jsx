import { useEffect, useState } from "react";
import Books from "./Books";

const AuthorProfile = ({ author }) => {
  const [books, setBooks] = useState([]);
  const getAuthorBooks = async () => {
    const response = await fetch("http://localhost:8080/getAuthorBooks", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ books: author.publishedBooks }),
    });
    if (response.status === 200) {
      const books = await response.json();
      setBooks(books);
    }
  };

  // Getting the book as soon as the profile is rendered.
  useEffect(() => {
    getAuthorBooks(); // calling the function to get the books
  }, []);

  return (
    <div className="fixed inset-0 flex m-14 justify-center z-50 font-inika">
      <div className="bg-white p-8 rounded-lg shadow-lg w-3/6 flex flex-col items-center">
        <img
          src={author.image_url}
          alt="Book"
          className="w-20 h-20 object-cover rounded-full"
        />
        <div className="font-bold text-2xl mt-3">{author.username}</div>
        <p className="text-subText text-sm">Sold: 600 Books</p>
        <p className="text-subText text-sm">
          Published Books: {author.publishedBooks.length}
        </p>
        <div className="w-2/3 mt-4 max-h-screen overflow-y-auto">
          {books === null ? <></> : <Books booksToShow={books} />}
        </div>
      </div>
    </div>
  );
};

export default AuthorProfile;
