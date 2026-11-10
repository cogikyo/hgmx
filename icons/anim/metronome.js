import { animate, τmax } from "../motion.js";

export default function (el) {
  const pendulum = el.querySelector("g");
  animate(pendulum, { rotate: [0, -8, 4, 0] }, { duration: τmax });
}
