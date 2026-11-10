import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const [cap, tassel] = el.querySelectorAll("g");
  animate(cap, { translateY: [0, -5, 1, 0], rotate: [0, -15, 5, 0] }, { duration: τmax });
  animate(tassel, { rotate: [0, 25, -10, 0] }, { duration: τmid, delay: τmin });
}
