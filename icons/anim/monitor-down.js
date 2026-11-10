import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
