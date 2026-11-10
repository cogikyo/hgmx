import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const post = el.querySelectorAll('path[d^="M10 22"], path[d^="M8 22"]');
  animate(post, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const sign = el.querySelector("g");
  animate(sign, { rotate: [0, -5, 5, -1.5, 0] }, { duration: τmax });
}
