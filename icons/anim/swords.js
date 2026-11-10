import { animate, τmax } from "../motion.js";

export default function (el) {
  const a = el.querySelector("g.a");
  animate(a, { rotate: [0, 10, -3, 0] }, { duration: τmax });
  const b = el.querySelector("g.b");
  animate(b, { rotate: [0, -10, 3, 0] }, { duration: τmax });
}
