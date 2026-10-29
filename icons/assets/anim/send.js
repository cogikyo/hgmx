import { animate, τmax } from "../motion.js";

export default function (el) {
  const plane = el.querySelector("g");
  animate(plane, { translateX: [0, 6, -6, 0], translateY: [0, -6, 6, 0], opacity: [1, 0, 0, 1] }, { duration: τmax, times: [0, 0.4, 0.41, 1] });
}
