import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector("g");
  animate(book, { translateY: [-3, 1, 0] }, { duration: τmid });

  const leaning = el.querySelector('path[d^="M20.4"]');
  animate(leaning, { rotate: [0, 12, 0] }, { duration: τmid, delay: τmin });
}
