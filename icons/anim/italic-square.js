import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const stem = el.querySelector('path[d^="m13.5 7"]');
  animate(stem, { rotate: [-20, 5, 0] }, { duration: τmid, delay: τmin });
  const top = el.querySelector('path[d^="M16 7"]');
  animate(top, { translateX: [-1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
  const bottom = el.querySelector('path[d^="M13 17"]');
  animate(bottom, { translateX: [1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
