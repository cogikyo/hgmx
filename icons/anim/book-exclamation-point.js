import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mark = el.querySelector("g");
  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d^="M12 13"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
