import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const up = el.querySelector('path[d^="m9 10"]');
  animate(up, { translateY: [0, -2, 1, 0] }, { duration: τmid, delay: τmin });

  const down = el.querySelector('path[d^="m15 14"]');
  animate(down, { translateY: [0, 2, -1, 0] }, { duration: τmid, delay: τmin });
}
