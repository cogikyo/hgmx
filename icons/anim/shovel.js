import { animate, τmax } from "../motion.js";

export default function (el) {
  const shovel = el.querySelectorAll("path");
  animate(shovel, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
