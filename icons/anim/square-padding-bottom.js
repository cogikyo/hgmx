import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bottom = el.querySelector('path[d="M8 17h8"]');
  animate(bottom, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
