import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M21 15"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const corner = el.querySelector('path[d^="M21 11"]');
  animate(corner, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
