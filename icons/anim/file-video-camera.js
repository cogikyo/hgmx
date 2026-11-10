import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 12V4"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector('path[d^="m10 17.843"]');
  animate(lens, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
