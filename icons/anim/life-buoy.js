import { animate, τmax } from "../motion.js";

export default function (el) {
  const buoy = el.querySelector("g");
  animate(buoy, { rotate: [0, -30, 12, 0], translateY: [0, 2, -1, 0] }, { duration: τmax });
}
