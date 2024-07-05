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
    <div className="font-inika flex flex-col w-full items-center pt-4 pb-4 pl-16 pr-16">
      <div className="font-bold text-2xl mb-4">Search Authors</div>
      <SearchBox onChange={searchAuthors} setAuthors={setAuthors} />
      <div className="mt-2">
        {authors.map((author) => (
          <div key={author.id}>{author.username}</div>
        ))}
      </div>
    </div>
  );
};

export default Search;
