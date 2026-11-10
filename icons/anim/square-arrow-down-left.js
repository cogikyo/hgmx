import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
