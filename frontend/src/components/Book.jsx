import { useDispatch, useSelector } from "react-redux";
import { incBoughtBook, setBook, deleteBook } from "state";
import { TbEdit } from "react-icons/tb";
import { MdDelete } from "react-icons/md";

const Book = ({ book, inAuthorProfile }) => {
  const dispatch = useDispatch();
  const publishDate = new Date(book.publish_date);
  const day = publishDate.getUTCDate();
  const month = publishDate.getUTCMonth() + 1;
  const year = publishDate.getUTCFullYear();
  const user = useSelector((state) => state.user);
  const isUserBook = user && user.id === book.publisher_id;

  const buyBook = async () => {
    // const bookId = book.id === undefined ? book._id : book.id;
    const bookId = book.id;

    try {
      const response = await fetch("http://localhost:8080/buyBook", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: bookId }),
      });
      if (response.status === 200) {
        dispatch(
          setBook({
            book: {
              ...book,
              quantity: book.quantity - 1,
            },
          })
        );
        if (inAuthorProfile === true) {
          book.quantity = book.quantity - 1;
        }
        dispatch(incBoughtBook());
      }
    } catch (err) {
      alert("Failed to buy book without internet connection");
    }
  };

  const deleteBookRequest = async () => {
    try {
      const response = await fetch("http://localhost:8080/deleteBook", {
        method: "DELETE",
        headers: { "content-type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ id: book.id }),
      });

      if (response.status === 200) {
        dispatch(deleteBook({ book: book }));
      }
    } catch {
      alert("could not delete book");
    }
  };
  const editBook = async () => {};

  return (
    <div className="flex flex-col border-black border-2 rounded-md w-full p-4 gap-2">
      <div className="flex flex-row gap-5 justify-center">
        <img
          src={book.image_url}
          alt="Book"
          className="w-52 h-62 object-cover rounded-lg"
        />
        <div className="flex flex-col w-full justify-evenly font-bold font-inika text-sm">
          <div className="text-2xl capitalize-first-letter">{book.title}</div>
          <div className="border border-gray-950" />
          <div className="font-inika">Author: {book.author}</div>
          <div>Left: {book.quantity}</div>
          <div>Publisher: {book.publisher}</div>
          <div>
            Publish date: {day}/{month}/{year}
          </div>
          <div className="flex flex-row gap-1">
            Price: <div className="text-buyMeButton">{book.price}$</div>
          </div>
        </div>
      </div>
      <div className="flex flex-row gap-5 items-center justify-center">
        <button
          className="bg-buyMeButton rounded-md p-3 text-white font-semibold flex-row flex justify-center items-center w-full"
          onClick={buyBook}
        >
          Buy now
        </button>
        {isUserBook && (
          <div className="flex flex-col justify-center items-center">
            <TbEdit
              size={22}
              className="hover:scale-110 transition ease-in-out duration-100 text-black"
              onClick={editBook}
            />
            <MdDelete
              size={22}
              className="hover:scale-110 transition ease-in-out duration-100 text-black"
              onClick={deleteBookRequest}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default Book;
