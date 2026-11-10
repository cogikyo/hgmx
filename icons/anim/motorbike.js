import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bike = el.querySelector("g");
  animate(bike, { translateX: [0, 3, -1, 0] }, { duration: τmax });
  animate(bike, { rotate: [0, -8, 0] }, { duration: τmid, delay: τmin });
}
