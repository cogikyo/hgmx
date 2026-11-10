import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hill = el.querySelector('path[d^="m20 13.7"]');
  animate(hill, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
