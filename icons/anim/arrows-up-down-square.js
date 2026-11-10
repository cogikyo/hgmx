import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const up = el.querySelector("g.up");
  animate(up, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmid });

  const down = el.querySelector("g.down");
  animate(down, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
