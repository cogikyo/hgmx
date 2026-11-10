import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const antenna = el.querySelector("g");
  animate(antenna, { rotate: [0, -15, 10, 0] }, { duration: τmax });

  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
