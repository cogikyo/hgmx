import { animate, τmax } from "../motion.js";

export default function (el) {
  const buoy = el.querySelector("g");
  animate(buoy, { rotate: [0, -15, 8, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
