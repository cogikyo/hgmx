import { animate, τmax } from "../motion.js";

export default function (el) {
  const bird = el.querySelector("g");
  animate(bird, { translateY: [0, -3, 0.5, 0] }, { duration: τmax });
}
