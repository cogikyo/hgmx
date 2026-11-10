import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M22 12"]');
  animate(head, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const eyes = el.querySelectorAll('path[d^="M10.7"], path[d^="M17.5"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
