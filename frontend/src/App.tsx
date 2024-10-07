import { useEffect, useState } from "react";

export default function App() {
  const [text, setText] = useState<string>();

  useEffect(() => {
    const fetchData = async () => {
      const res = await fetch("/api/");

      if (res.status === 200) {
        setText(await res.text());
      }
    };

    fetchData();
  }, []);

  return <div>{text}</div>;
}
