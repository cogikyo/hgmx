import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M14.692"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const box = el.querySelector("g");
  animate(box, { translateY: [0, -3, 1, 0], scaleY: [1, 1.1, 0.9, 1] }, { duration: τmid, delay: τmin });
}
