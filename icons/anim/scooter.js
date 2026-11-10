import { animate, τmax } from "../motion.js";

export default function (el) {
  const scooter = el.querySelectorAll("path, circle");
  animate(scooter, { translateX: [0, 3, -1, 0] }, { duration: τmax });
}
