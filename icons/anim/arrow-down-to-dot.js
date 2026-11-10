import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [-6, 2, 0], opacity: [0, 1] }, { duration: τmax });

  const dot = el.querySelector("circle");
  animate(dot, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
