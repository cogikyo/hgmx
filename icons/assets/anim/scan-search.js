import { animate, τmax } from "../motion.js";

export default function (el) {
  const [corners, lens] = el.querySelectorAll("g");
  animate(corners, { scale: [1.15, 0.95, 1] }, { duration: τmax });

  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, 4, -2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
}
