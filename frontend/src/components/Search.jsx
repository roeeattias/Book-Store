import { useState } from "react";
import SearchBox from "./SearchBox";

const Search = () => {
  const [authors, setAuthors] = useState([]);
  const searchAuthors = async (query) => {
    const formData = new FormData();
    formData.append("query", query);

    const response = await fetch("http://localhost:8080/getAuthors", {
      method: "POST",
      body: formData,
    });

    if (response.status === 200) {
      const authors = await response.json();
      setAuthors(authors);
    }
  };

  return (
    <div className="font-inika flex flex-col w-full items-center pb-4 pl-16 pr-16">
      <div className="font-bold text-2xl mb-4">Search Authors</div>
      <SearchBox onChange={searchAuthors} setAuthors={setAuthors} />
      <div className="flex flex-col mt-2 w-full">
        {authors.map((author) => (
          <div
            key={author.id}
            className="flex flex-row w-full rounded-lg border border-black border-opacity-40 p-2 mb-2 items-center gap-5 text-lg hover:bg-gray-300"
            onClick={() => {
              console.log(author.username);
            }}
          >
            <img
              src={`${process.env.PUBLIC_URL}/images/author_image.jpeg`}
              alt="Book"
              className="w-12 h-12 object-cover rounded-full"
            />
            {author.username}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Search;
