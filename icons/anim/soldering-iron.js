import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const iron = el.querySelector("g");
  animate(iron, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmid });

  const smoke = el.querySelector('path[d^="M3 2"]');
  animate(smoke, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
