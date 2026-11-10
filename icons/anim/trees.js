import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const [round, pine] = el.querySelectorAll("g");
  animate(round, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
  animate(pine, { rotate: [0, -14, 8, 0] }, { duration: τmid, delay: τmin });
}
