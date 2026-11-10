import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M11 19"], path[d^="M13 5"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector("circle");
  animate(lens, { scaleX: [1, 0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
