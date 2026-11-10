import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const top = el.querySelector('path[d="M8 7h8"]');
  animate(top, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
