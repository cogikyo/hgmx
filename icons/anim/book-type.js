import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const letter = el.querySelectorAll('path[d^="M16 8"], path[d="M12 6v7"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const foot = el.querySelector('path[d="M10 13h4"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
