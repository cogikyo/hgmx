import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M11.65"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 8, 0] }, { duration: τmid, delay: τmin });
}
