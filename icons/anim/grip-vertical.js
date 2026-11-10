import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const dots = el.querySelectorAll("circle");
  animate(dots, { scale: [1, 0.5, 1] }, { duration: τmid });

  const grip = el.querySelector("g");
  animate(grip, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
