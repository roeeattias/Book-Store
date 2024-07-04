const Book = ({ book }) => {
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
          <div className="text-red-600">Publish date: // need to parse</div>
          <div className="text-red-600">Price: // need to add</div>
        </div>
      </div>
      <button className="bg-buyMeButton m-2 rounded-md p-3 text-white font-semibold">
        Buy now
      </button>
    </div>
  );
};

export default Book;
