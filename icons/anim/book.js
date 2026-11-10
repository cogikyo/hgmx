import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const book = el.querySelector("path");
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(book, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
