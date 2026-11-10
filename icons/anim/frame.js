import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rails = el.querySelectorAll('line[x1="22"]');
  animate(rails, { scaleX: [0, 1.1, 1] }, { duration: τmid });

  const stiles = el.querySelectorAll('line[y1="2"]');
  animate(stiles, { scaleY: [0, 1.1, 1] }, { duration: τmid, delay: τmin });
}
