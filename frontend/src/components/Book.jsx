import { useDispatch } from "react-redux";
import { incBoughtBook, setBook } from "state";

const Book = ({ book }) => {
  const dispatch = useDispatch();
  const publishDate = new Date(book.publish_date);
  const day = publishDate.getUTCDate();
  const month = publishDate.getUTCMonth() + 1;
  const year = publishDate.getUTCFullYear();

  const buyBook = async () => {
    const response = await fetch("http://localhost:8080/buyBook", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: book.id }),
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
      dispatch(incBoughtBook());
    }
  };

  return (
    <div className="flex flex-col border-black border-2 rounded-md w-full p-3">
      <div className="flex flex-row m-2 gap-6">
        <img
          src={`${process.env.PUBLIC_URL}/images/book_image.jpg`}
          alt="Book"
          className="w-64 h-64 object-cover rounded-lg"
        />
        <div className="flex flex-col w-full justify-evenly font-bold font-inika">
          <div className="text-3xl">{book.title}</div>
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
      <button
        className="bg-buyMeButton m-2 rounded-md p-3 text-white font-semibold"
        onClick={buyBook}
      >
        Buy now
      </button>
    </div>
  );
};

export default Book;
