import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bar = el.querySelector('path[d^="M7 7"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const legs = el.querySelector("g");
  animate(legs, { scaleY: [0, 1] }, { duration: τmin, delay: τmid });
}
