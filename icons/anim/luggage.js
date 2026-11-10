import { animate, τmax } from "../motion.js";

export default function (el) {
  const luggage = el.querySelector("g");
  animate(luggage, { rotate: [0, -8, 3, 0] }, { duration: τmax });
}
