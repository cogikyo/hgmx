import { animate, τmax } from "../motion.js";

export default function (el) {
  const knob = el.querySelector("circle");
  animate(knob, { scaleX: [1, 1.3, 0.95, 1] }, { duration: τmax });
}
